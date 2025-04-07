import { useEffect, type ReactNode } from 'react';
import useBaseUrl from '@docusaurus/useBaseUrl';
import Layout from '@theme/Layout';

export default function Home(): ReactNode {
  const redirectUrl = useBaseUrl('/docs/getting-started');
  useEffect(() => {
    window.location.href = redirectUrl;
  }, [redirectUrl]);
  return (
    <Layout>
      <main></main>
    </Layout>
  );
}
