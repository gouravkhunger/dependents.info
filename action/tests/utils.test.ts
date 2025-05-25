import { validateRepoName } from "@/utils";

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
