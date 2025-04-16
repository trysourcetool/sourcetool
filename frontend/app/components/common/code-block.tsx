import { createHighlighter } from 'shiki';
import { useEffect, useRef, useState, type FC } from 'react';

export const CodeBlock: FC<{
  code: string;
  language: string;
}> = ({ code, language }) => {
  const isInitialLoading = useRef(false);

  const [formattedCode, setFormattedCode] = useState('');

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;
    (async () => {
      const highlighter = await createHighlighter({
        themes: ['github-dark-dimmed'],
        langs: [language],
      });

      const html = highlighter.codeToHtml(code, {
        lang: language,
        theme: 'github-dark-dimmed',
      });

      setFormattedCode(html);
      isInitialLoading.current = false;
    })();
  }, [code, language]);

  return (
    <div
      className="overflow-hidden rounded-md text-xs font-normal [&>pre]:px-4 [&>pre]:py-3"
      dangerouslySetInnerHTML={{ __html: formattedCode }}
    />
  );
};
