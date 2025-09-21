export const getLevelColor = (level: string) => {
  switch (level.toLowerCase()) {
    case 'error': return 'text-red-400';
    case 'warn': case 'warning': return 'text-yellow-400';
    case 'info': return 'text-blue-400';
    case 'debug': return 'text-purple-400';
    case 'success': return 'text-green-400';
    default: return 'text-gray-400';
  }
};

export const getLevelBg = (level: string) => {
  switch (level.toLowerCase()) {
    case 'error': return 'bg-red-500/10 border-red-500/20';
    case 'warn': case 'warning': return 'bg-yellow-500/10 border-yellow-500/20';
    case 'info': return 'bg-blue-500/10 border-blue-500/20';
    case 'debug': return 'bg-purple-500/10 border-purple-500/20';
    case 'success': return 'bg-green-500/10 border-green-500/20';
    default: return 'bg-gray-500/10 border-gray-500/20';
  }
};

export const formatTime = (time: string) => {
  try {
    return new Date(time).toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  } catch {
    return time;
  }
};

export const truncateMessage = (message: string, maxLength: number = 100) => {
  if (message.length <= maxLength) return message;
  return message.substring(0, maxLength) + '...';
};