import { Sourcetool, SourcetoolConfig, UIBuilder } from '@sourcetool/node';
import { listUsers } from './users';

const listUsersPage = async (ui: UIBuilder) => {
  const searchCols = ui.columns(2);
  const name = searchCols[0].textInput('Name', {
    placeholder: 'Enter name to filter',
  });
  const email = searchCols[1].textInput('Email', {
    placeholder: 'Enter email to filter',
  });

  const users = listUsers(name, email, 0, '');

  ui.table(users, {
    height: 10,
    columnOrder: ['ID', 'Name', 'Email', 'Age', 'Gender', 'CreatedAt'],
  });
};

const config: SourcetoolConfig = {
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/users', 'Users', listUsersPage);

sourcetool.listen();
