/*
 * Copyright 2018, Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hipstershop;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import hipstershop.Demo.Ad;
import hipstershop.Demo.AdRequest;
import hipstershop.Demo.AdResponse;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.util.Map;
import java.util.HashMap;
import java.util.Base64;
import java.util.logging.Level;
import java.util.logging.Logger;

public final class AdService {

  private static final Logger logger = Logger.getLogger(AdService.class.getName());

  @SuppressWarnings("FieldCanBeLocal")
  private static int MAX_ADS_TO_SERVE = 2;

  private static final ImmutableListMultimap<String, Ad> adsMap = createAdsMap();

  private static class RequestContext {
    public HTTP http = new HTTP();

    static class HTTP {
      public String method;
      public String path;
    }
  }

  private static class RequestData {
    public String body;
    public Map<String, String> headers;
    public RequestContext requestContext = new RequestContext();
    public boolean isBase64Encoded;
  }

  private static class ResponseData {
    public int statusCode;
    public Map<String, String> headers;
    public String body;
    public boolean isBase64Encoded;
  }

  private Collection<Ad> getAdsByCategory(String category) {
    return adsMap.get(category);
  }

  private List<Ad> getRandomAds() {
    List<Ad> ads = new ArrayList<>(MAX_ADS_TO_SERVE);
    Collection<Ad> allAds = adsMap.values();
    for (int i = 0; i < MAX_ADS_TO_SERVE; i++) {
      ads.add(Iterables.get(allAds, random.nextInt(allAds.size())));
    }
    return ads;
  }

  private static ImmutableListMultimap<String, Ad> createAdsMap() {
    Ad hairdryer =
        Ad.newBuilder()
            .setRedirectUrl("/product/2ZYFJ3GM2N")
            .setText("Hairdryer for sale. 50% off.")
            .build();
    Ad tankTop =
        Ad.newBuilder()
            .setRedirectUrl("/product/66VCHSJNUP")
            .setText("Tank top for sale. 20% off.")
            .build();
    Ad candleHolder =
        Ad.newBuilder()
            .setRedirectUrl("/product/0PUK6V6EV0")
            .setText("Candle holder for sale. 30% off.")
            .build();
    Ad bambooGlassJar =
        Ad.newBuilder()
            .setRedirectUrl("/product/9SIQT8TOJO")
            .setText("Bamboo glass jar for sale. 10% off.")
            .build();
    Ad watch =
        Ad.newBuilder()
            .setRedirectUrl("/product/1YMWWN1N4O")
            .setText("Watch for sale. Buy one, get second kit for free")
            .build();
    Ad mug =
        Ad.newBuilder()
            .setRedirectUrl("/product/6E92ZMYYFZ")
            .setText("Mug for sale. Buy two, get third one for free")
            .build();
    Ad loafers =
        Ad.newBuilder()
            .setRedirectUrl("/product/L9ECAV7KIM")
            .setText("Loafers for sale. Buy one, get second one for free")
            .build();
    return ImmutableListMultimap.<String, Ad>builder()
        .putAll("clothing", tankTop)
        .putAll("accessories", watch)
        .putAll("footwear", loafers)
        .putAll("hair", hairdryer)
        .putAll("decor", candleHolder)
        .putAll("kitchen", bambooGlassJar, mug)
        .build();
  }

  private static AdRequest decodeRequest(String request) throws Exception {
    RequestData reqData = JsonFormat.parser().merge(request, RequestData.class);
    byte[] binReqBody;

    if (reqData.isBase64Encoded) {
      binReqBody = Base64.getDecoder().decode(reqData.body);
    } else {
      binReqBody = reqData.body.getBytes();
    }

    AdRequest.Builder adRequestBuilder = AdRequest.newBuilder();
    adRequestBuilder.mergeFrom(binReqBody);
    return adRequestBuilder.build();
  }

  private static AdResponse handleRequest(AdRequest req) {
    List<Ad> allAds = new ArrayList<>();
    logger.info("Received ad request (context_words=" + req.getContextKeysList() + ")");

    if (req.getContextKeysCount() > 0) {
      for (int i = 0; i < req.getContextKeysCount(); i++) {
        Collection<Ad> ads = service.getAdsByCategory(req.getContextKeys(i));
        allAds.addAll(ads);
      }
    } else {
      allAds = getRandomAds();
    }

    if (allAds.isEmpty()) {
      // Serve random ads.
      allAds = getRandomAds();
    }

    return AdResponse.newBuilder().addAllAds(allAds).build();
  }

  private static String encodeResponse(AdResponse adResponse) throws InvalidProtocolBufferException {
    byte[] binRespBody = adResponse.toByteArray();
    String encodedRespBody = Base64.getEncoder().encodeToString(binRespBody);

    ResponseData respData = new ResponseData();
    respData.statusCode = 200;
    respData.headers = new HashMap<>();
    respData.headers.put("Content-Type", "application/octet-stream");
    respData.body = encodedRespBody;
    respData.isBase64Encoded = true;

    return JsonFormat.printer().print(respData);
  }

  public static void main(String[] args) {
    try {
      BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
      String request = reader.readLine();

      AdRequest adRequest = decodeRequest(request);
      AdResponse adResponse = handleRequest(adRequest);
      String response = encodeResponse(adResponse);

      System.out.println(response);
    } catch (Exception e) {
      logger.log(Level.SEVERE, "Error occurred", e);
      System.exit(1);
    }
  }
}
