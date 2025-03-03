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
    (async () => {
      const highlighter = await createHighlighter({
        themes: ['github-dark-dimmed'],
        langs: [language],
      });

      console.log({ code });

      const html = highlighter.codeToHtml(code, {
        lang: language,
        theme: 'github-dark-dimmed',
      });

      console.log({ html });

      setFormattedCode(html);
      isInitialLoading.current = true;
    })();
  }, [code, language]);

  return (
    <div
      className="overflow-hidden rounded-md text-xs font-normal [&>pre]:px-4 [&>pre]:py-3"
      dangerouslySetInnerHTML={{ __html: formattedCode }}
    />
  );
};
