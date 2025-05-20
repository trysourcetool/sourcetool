import {
  TableOnSelect,
  Sourcetool,
  SourcetoolConfig,
  UIBuilder,
} from '@sourcetool/node';
import { listUsers, RefundRequest, refundStripe } from './user.ts';

const refundPage = async (ui: UIBuilder) => {
  const searchCols = ui.columns(2);
  const name = searchCols[0].textInput('Name', {
    placeholder: 'Enter user name to filter',
  });
  const email = searchCols[1].textInput('Email', {
    placeholder: 'Enter email to filter',
  });

  const users = listUsers(name, email);

  const tableWidget = ui.table(users, {
    height: 10,
    columnOrder: ['ID', 'Name', 'Email', 'CreatedAt'],
    onSelect: TableOnSelect.Rerun,
  });

  const selectedUser =
    tableWidget?.selection?.row && tableWidget.selection.row < users.length
      ? users[tableWidget.selection.row]
      : null;

  if (selectedUser) {
    const [formWidget, submitted] = ui.form('Refund', { clearOnSubmit: true });
    const amount = formWidget.numberInput('Amount', {
      minValue: 1,
      required: true,
    });
    const reason = formWidget.textInput('Reason', {
      placeholder: 'Enter refund reason',
      required: true,
    });

    if (submitted) {
      const refundReq: RefundRequest = {
        userId: selectedUser.id,
        amount: amount ?? 0,
        reason: reason,
      };
      refundStripe(refundReq);
      ui.markdown(
        `Refund processed for user ${selectedUser.name} (${selectedUser.email}), amount: ${refundReq.amount}`,
      );
    }
  }
};

const config: SourcetoolConfig = {
  apiKey: 'your_api_key',
  endpoint: 'ws://localhost:3000',
};

const sourcetool = new Sourcetool(config);

sourcetool.page('/refunds', 'Refunds', refundPage);

sourcetool.listen();
