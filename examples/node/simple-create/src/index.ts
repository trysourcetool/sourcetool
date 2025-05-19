import { Sourcetool, SourcetoolConfig, UIBuilderType } from '@sourcetool/node';
import { createUser } from './users';

const createUserPage = async (ui: UIBuilderType) => {
  const [form, submitted] = ui.form('Create User', {
    clearOnSubmit: true,
  });

  const formName = form.textInput('Name', {
    placeholder: 'Enter user name',
    required: true,
  });

  const formEmail = form.textInput('Email', {
    placeholder: 'Enter user email',
    required: true,
  });

  const formAge = form.numberInput('Age', {
    minValue: 0,
    maxValue: 100,
  });

  const formGender = form.selectbox('Gender', {
    options: ['male', 'female'],
  });

  if (submitted) {
    const user = createUser({
      id: '30',
      name: formName,
      email: formEmail,
      age: formAge || 0,
      gender: formGender?.value ?? '',
      createdAt: new Date(),
    });

    if (user) {
      ui.markdown(`User created successfully with ID: ${user.id}`);
    }
  }
};

const config: SourcetoolConfig = {
  apiKey: 'development_ZhMe3YBmFMSeE9Vx9DcsVQ1E6WIZvBIUZhMe3YBmFMSeE9Vx9Dc',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/users/new', 'Create User', createUserPage);

sourcetool.listen();
