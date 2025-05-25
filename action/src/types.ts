export type Dependents = {
  name: string;
  stars: number;
  image: string;
}[];

export type DependentsPage = {
  nextPageLink?: string;
  dependents: Dependents;
};
