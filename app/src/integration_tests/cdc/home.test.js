import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { runContract } from './cdc';

runContract("home/no_conversations.yaml", (_) => {
  expect(screen.getByText(/Explore Art/i)).toBeInTheDocument()
});

runContract("home/multiple_conversations.yaml", (app, _) => {
  expect(screen.getByText(/Explore Art/i)).toBeInTheDocument()

  app.textShouldExist.forEach((text, _) => {
    expect(screen.getByText(text)).toBeInTheDocument()
  });
});

runContract("home/start_conversation.yaml",
  (_) => {
    expect(screen.getByText(/Tell me something about:/i)).toBeInTheDocument()
  },
  () => {
    userEvent.click(screen.getByText(/Explore Art/i));
  }
);
