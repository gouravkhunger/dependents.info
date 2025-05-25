export const validateRepoName = (repoName: string): boolean => {
  const regex = /^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/;
  return regex.test(repoName);
};
