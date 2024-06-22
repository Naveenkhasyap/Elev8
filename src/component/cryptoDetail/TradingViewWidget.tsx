import React, { useEffect, useRef } from "react";

declare global {
  interface Window {
    TradingView: any;
  }
}

interface TradingViewWidgetProps {
  symbol: string;
}

const TradingViewWidget: React.FC<TradingViewWidgetProps> = ({ symbol }) => {
  const widgetRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (widgetRef.current) {
      const script = document.createElement("script");
      script.src =
        "../../../charting_library-master/charting_library/charting_library.standalone.js";
      // script.src =
      //   "https://s3.tradingview.com/external-embedding/embed-widget-advanced-chart.js";
      script.async = true;
      script.onload = () => {
        console.log("TradingView script loaded");
        new window.TradingView.widget({
          container_id: widgetRef.current!,
          width: "100%",
          height: "100%",
          // symbol: symbol,
          symbol: "AAPL",
          locale: "en",
          fullscreen: true,
          interval: "D",
          timezone: "Etc/UTC",
          library_path: "../../../charting_library-master/charting_library",
          theme: "dark",
          style: "1",
          toolbar_bg: "#f1f3f6",
          enable_publishing: false,
          allow_symbol_change: true,
          hide_side_toolbar: false,
          details: true,
          hotlist: true,
          calendar: true,
        });
      };
      script.onerror = () => {
        console.error("Error loading TradingView script");
      };
      widgetRef.current.appendChild(script);

      return () => {
        widgetRef.current!.innerHTML = "";
      };
    }
  }, [symbol]);

  return (
    <div
      ref={widgetRef}
      style={{ border: "1px solid red " }}
      className="w-full h-screen"
    />
  );
};

export default TradingViewWidget;
