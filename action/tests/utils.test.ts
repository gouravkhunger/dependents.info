import { imageUrlToBase64, validateRepoName } from "@/utils";

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
});

describe("image URL to Base64 conversion", () => {
  it("should convert image URL to Base64 string", async () => {
    const url = "https://gourav.sh/assets/images/gourav.webp";
    const data = await imageUrlToBase64(url);
    console.log(data);
    expect(data).toMatch(/^data:image\/webp;base64,/);
  });

  it("should handle non-image URLs gracefully", async () => {
    const url = "https://example.com/not-an-image.txt";
    return expect(imageUrlToBase64(url)).rejects.toThrow();
  });
});
