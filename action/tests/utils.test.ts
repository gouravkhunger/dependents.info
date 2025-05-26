import { imageUrlToBase64, removeQueryParams, validateRepoName } from "@/utils";

describe("utility functions", () => {
  it("validateRepoName should validate github repo names correctly", () => {
    expect(validateRepoName("")).toBe(false);
    expect(validateRepoName("  ")).toBe(false);
    expect(validateRepoName("user repo")).toBe(false);
    expect(validateRepoName(" user/repo ")).toBe(false);
    expect(validateRepoName("organization/repo/subdir")).toBe(false);

    expect(validateRepoName("user/repo")).toBe(true);
    expect(validateRepoName("organization/repo")).toBe(true);
  });

  it("removeQueryParams should remove specified query parameters from URL", () => {
    const invalidUrl = "not-a-valid-url";
    const url = "https://example.com/path?param1=value1&param2=value2";
    expect(() => removeQueryParams(invalidUrl, "param")).toThrow();
    expect(removeQueryParams(url, "param1")).toBe(
      "https://example.com/path?param2=value2",
    );
    expect(removeQueryParams(url, "param2")).toBe(
      "https://example.com/path?param1=value1",
    );
    expect(removeQueryParams(url, "param1", "param2")).toBe(
      "https://example.com/path",
    );
    expect(removeQueryParams(url, "nonexistent")).toBe(url);
  });
});

describe("image URL to Base64 conversion", () => {
  it("should convert image URL to Base64 string", async () => {
    const url = "https://gourav.sh/assets/images/gourav.webp";
    const data = await imageUrlToBase64(url);
    expect(data).toMatch(/^data:image\/webp;base64,/);
  });

  it("should handle non-image URLs gracefully", async () => {
    const url = "https://gourav.sh/404";
    return expect(imageUrlToBase64(url)).rejects.toThrow();
  });
});
