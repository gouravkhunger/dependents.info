import { API_BASE_URL, ERROR } from "@/constants";
import {
  buildAPIUrl,
  buildDependentsUrl,
  imageUrlToBase64,
  removeQueryParams,
  validateRepoName,
} from "@/utils";

describe("utility functions", () => {
  it("buildAPIUrl should build the correct GitHub dependents URL", () => {
    expect(() => buildAPIUrl("invalid/repo/name", "")).toThrow(
      ERROR.INVALID_REPO_FORMAT,
    );
    expect(buildAPIUrl("owner/repo", "")).toBe(
      `${API_BASE_URL}/owner/repo/ingest`,
    );
    expect(buildAPIUrl("owner/repo", "random_id")).toBe(
      `${API_BASE_URL}/owner/repo/ingest?id=random_id`,
    );
  });

  it("buildDependentsUrl should build the correct GitHub dependents URL", () => {
    expect(() => buildDependentsUrl("invalid/repo/name", "")).toThrow(
      ERROR.INVALID_REPO_FORMAT,
    );
    expect(buildDependentsUrl("owner/repo", "")).toBe(
      "https://github.com/owner/repo/network/dependents",
    );
    expect(buildDependentsUrl("owner/repo", "random_id")).toBe(
      "https://github.com/owner/repo/network/dependents?package_id=random_id",
    );
  });

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
