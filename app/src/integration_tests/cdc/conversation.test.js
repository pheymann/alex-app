import { screen, fireEvent } from '@testing-library/react';
import { runContract } from './cdc';
import userEvent from '@testing-library/user-event';

runContract("conversation/found.yaml", (app) => {
  expect(screen.getByPlaceholderText(/Do you have any questions?/i)).toBeInTheDocument()

  app.textShouldExist.forEach((text, _) => {
    expect(screen.getByText(text)).toBeInTheDocument()
  });
});

runContract("conversation/not_found.yaml", (_, counter) => {
  expect(screen.getByText(/Explore Art/i)).toBeInTheDocument();
  expect(counter.get('/api/app/logs')).toBe(1);
});

runContract("conversation/unauthorized_access.yaml", (_, counter) => {
  expect(screen.getByText(/Explore Art/i)).toBeInTheDocument();
  expect(counter.get('/api/app/logs')).toBe(1);
});

runContract("conversation/start_conversation.yaml",
  (app, _) => {
    app.textShouldExist.forEach((text, _) => {
      expect(screen.getByText(text)).toBeInTheDocument()
    });
  },
  () => {
    fireEvent.change(screen.getByPlaceholderText(/The Mona Lisa by Leonardo da Vinci/i), { target: { value: 'art' } });
    userEvent.click(screen.getByText(/Send/i));
  },
);

runContract("conversation/ask_question.yaml",
  (app, _) => {
    app.textShouldExist.forEach((text, _) => {
      expect(screen.getByText(text)).toBeInTheDocument()
    });
  },
  () => {
    fireEvent.change(screen.getByPlaceholderText(/Do you have any questions?/i), { target: { value: 'another question' } });
    userEvent.click(screen.getByText(/Send/i));
  },
);
