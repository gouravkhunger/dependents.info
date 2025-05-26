export type DependentsPage = {
  nextPageLink?: string;
  dependents: Dependents;
};

export type Dependents = {
  name: string;
  stars: number;
  image: string;
}[];
