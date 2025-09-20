import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { CancelKumo, RunKumo } from '@wailsjs/go/main/App';
import { main } from '@wailsjs/go/models';
import { Loader2Icon } from 'lucide-react';

interface RunButtonProps {
  isBuilding: boolean;
  setIsBuilding: (isBuilding: boolean) => void;
  config: any;
  className?: string;
}

const RunButton = ({ isBuilding, setIsBuilding, config, className }: RunButtonProps) => {
  const handleToggle = async () => {
    if (isBuilding) {
      await CancelKumo();
      setIsBuilding(false);
    } else {
      setIsBuilding(true);
      try {
        await RunKumo(new main.KumoConfig(config));
      } catch (err) {
        console.error("Error running Kumo:", err);
      } finally {
        setIsBuilding(false);
      }
    }
  };

  return (
    <Button
      onClick={handleToggle}
      className={className}
    >
      {isBuilding && <Loader2Icon className="animate-spin" />}
      {isBuilding ? 'Stop' : 'Run'}
    </Button>
  );
};

export default RunButton;