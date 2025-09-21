import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import React from 'react';

interface LogsFooterProps {
  autoScroll: boolean;
  handleAutoScrollChange: (checked: boolean) => void;
  filteredCount: number;
  totalCount: number;
}

const LogsFooter: React.FC<LogsFooterProps> = ({
  autoScroll,
  handleAutoScrollChange,
  filteredCount,
  totalCount,
}) => (
  <>
    <div className="flex items-center gap-4">
      <Label className="flex items-center gap-2 text-sm">
        <Checkbox
          checked={autoScroll}
          onCheckedChange={(check) => handleAutoScrollChange(!!check)}
        />
        Auto-scroll to bottom
      </Label>
    </div>
    <div className="text-xs text-secondary-foreground">
      Showing {filteredCount} of {totalCount} logs
    </div>
  </>
);

export default LogsFooter;