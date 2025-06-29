import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 10 },
    { duration: "3m", target: 10 },
    { duration: "30s", target: 0 },
  ],
  tags: {
    name: "workshop-01",
    type: "load-test",
    testid: "workshop-01",
  },
//   scenarios: {
//     get_product_information: {
//       executor: "constant-vus",
//       vus: 10,
//       duration: "3m",
//     },
//   },
};

export default function () {
  // Get all users
  let res = http.get("http://192.168.3.140/api/v1/product/1");
  check(res, {
    "status is 200": (r) => r.status === 200,
    "body is not empty": (r) => r.body.length > 0,
  });

  // Validate with jsonschema
  const schema = {
    type: "object",
    required: [
      "id",
      "product_name",
      "product_price",
      "product_price_thb",
      "product_price_full_thb",
      "product_image",
      "stock",
      "product_brand",
    ],
    properties: {
      id: { type: "number" },
      product_name: { type: "string" },
      product_price: { type: "number" },
      product_price_thb: { type: "number" },
      product_price_full_thb: { type: "number" },
      product_image: { type: "string" },
      stock: { type: "number" },
      product_brand: { type: "string" },
    },
  };

  check(res, {
    "response matches schema": (r) => r.json(schema),
  });
}
