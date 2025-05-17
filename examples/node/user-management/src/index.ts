import {
  TableOnSelect,
  Sourcetool,
  SourcetoolConfig,
  UIBuilder,
} from '@sourcetool/node';
import { createUser, listUsers } from './users';

const listUsersPage = async (ui: UIBuilder) => {
  const searchCols = ui.columns(2);
  const name = searchCols[0].textInput('Name', {
    placeholder: 'Enter name to filter',
  });
  const email = searchCols[1].textInput('Email', {
    placeholder: 'Enter email to filter',
  });

  const users = listUsers(name, email, 0, '');

  const baseCols = ui.columns(2, { weight: [3, 1] });

  const table = baseCols[0].table(users, {
    height: 10,
    columnOrder: ['ID', 'Name', 'Email', 'Age', 'Gender', 'CreatedAt'],
    onSelect: TableOnSelect.Rerun,
  });

  let defaultName: string = '';
  let defaultEmail: string = '';
  let defaultGender: string = 'male';
  let defaultAge: number = 0;
  if (table?.selection && table.selection.row < users.length) {
    const selectedData = users[table.selection.row];

    if (selectedData) {
      defaultName = selectedData.name;
      defaultEmail = selectedData.email;
      defaultAge = selectedData.age;
      defaultGender = selectedData.gender;
    }
  }

  const [form, submitted] = baseCols[1].form('Update', {
    clearOnSubmit: true,
  });
  const formName = form.textInput('Name', {
    placeholder: 'Enter your name',
    defaultValue: defaultName,
    required: true,
  });

  const formEmail = form.textInput('Email', {
    placeholder: 'Enter your email',
    defaultValue: defaultEmail,
  });

  const formAge = form.numberInput('Age', {
    minValue: 0,
    maxValue: 100,
    defaultValue: defaultAge,
  });

  const formGender = form.selectbox('Gender', {
    options: ['male', 'female'],
    defaultValue: defaultGender,
  });

  if (submitted) {
    const user = createUser({
      id: '30',
      name: formName,
      email: formEmail,
      age: formAge ?? 0,
      gender: formGender?.value ?? '',
      createdAt: new Date(),
    });

    if (user) {
      ui.markdown(`User created successfully with ID: ${user.id}`);
    }
  }
};

const createUserPage = async (ui: UIBuilder) => {
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
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/users', 'Users', listUsersPage);
sourcetool.page('/users/new', 'Create User', createUserPage);

sourcetool.listen();
