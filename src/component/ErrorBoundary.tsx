import React, { ReactNode, useEffect } from "react";
import { ErrorBoundary } from "react-error-boundary";

const ChartErrorBoundary: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  useEffect(() => {
    let chartInstance: Chart | null = null;

    return () => {
      if (chartInstance) {
        chartInstance.destroy();
        chartInstance = null;
      }
    };
  }, []);

  return (
    <ErrorBoundary fallback={<div>Something went wrong with the chart.</div>}>
      {children}
    </ErrorBoundary>
  );
};
export default ChartErrorBoundary;
