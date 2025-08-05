import { useEffect } from 'react';
import { createFileRoute, useNavigate } from '@tanstack/react-router';

export default function AgentDetail() {
  const { agentId } = Route.useParams();
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect to the first chat or create a new one
    const defaultChatId = '1'; // Default to first chat
    navigate({
      to: '/agents/$agentId/chat/$chatId',
      params: { agentId, chatId: defaultChatId },
      replace: true,
    });
  }, [agentId, navigate]);

  return null; // This component just redirects
}

export const Route = createFileRoute('/_auth/agents/$agentId')({
  component: AgentDetail,
});