import {
  createContext,
  useContext,
  useState,
  type Dispatch,
  type ReactNode,
  type SetStateAction,
} from 'react';

type Breadcrumb = {
  to?: string;
  label: string;
};

type BreadcrumbsState = {
  breadcrumbsState?: Breadcrumb[];
  setBreadcrumbsState?: Dispatch<SetStateAction<Breadcrumb[]>>;
};

export const breadcrumbsContext = createContext<BreadcrumbsState>({});

export function BreadcrumbsProvider(props: { children: ReactNode }) {
  const [breadcrumbsState, setBreadcrumbsState] = useState<Breadcrumb[]>([]);

  return (
    <breadcrumbsContext.Provider
      value={{ breadcrumbsState, setBreadcrumbsState }}
    >
      {props.children}
    </breadcrumbsContext.Provider>
  );
}

export const useBreadcrumbs = () => {
  const { breadcrumbsState, setBreadcrumbsState } =
    useContext(breadcrumbsContext);

  return { breadcrumbsState, setBreadcrumbsState };
};
