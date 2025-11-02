import { fetch } from "./instance";

export async function Verify() {
  try {
    const res = await fetch({
      url: "/api/v1/auth/verify",
      method: "POST",
    });

    return res.data;
  } catch (err) {
    return { auth: false };
  }
}
