import { Sourcetool, SourcetoolConfig, UIBuilderType } from '@sourcetool/node';

const helloPage = async (ui: UIBuilderType) => {
  ui.markdown('# Hello, Sourcetool!');
  ui.markdown(
    'This is a simple example demonstrating the basic usage of the Sourcetool Go SDK.',
  );

  const name = ui.textInput('Your Name', {
    placeholder: 'Enter your name',
  });

  if (name !== '') {
    ui.markdown(`## Hello, ${name}!`);
    ui.markdown('Welcome to Sourcetool!');
  }

  const clicked = ui.button('Say Hello', {
    disabled: false,
  });

  if (clicked) {
    ui.markdown('👋 Hello from the button click!');
  }
};

const config: SourcetoolConfig = {
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/hello', 'Hello', helloPage);

sourcetool.listen();
