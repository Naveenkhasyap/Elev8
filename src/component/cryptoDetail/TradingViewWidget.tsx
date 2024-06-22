import React, { useEffect, useRef, useState } from "react";
import Script from "next/script";

interface CustomWindow extends Window {
  TradingView: any;
  Datafeeds: any;
}

const TradingViewWidget = () => {
  const widgetRef = useRef<HTMLDivElement>(null);
  const [scriptsLoaded, setScriptsLoaded] = useState(false);

  useEffect(() => {
    if (!scriptsLoaded) return;

    if (widgetRef.current) {
      (window as CustomWindow).TradingView.widget({
        container_id: widgetRef.current!,
        locale: "en",
        library_path: "/charting_library-master/charting_library/",
        datafeed: new (window as CustomWindow).Datafeeds.UDFCompatibleDatafeed(
          "https://demo-feed-data.tradingview.com"
        ),
        symbol: "AAPL",
        interval: "1D",
        fullscreen: true,
        debug: true,
      });
    }
  }, [scriptsLoaded]);

  return (
    <>
      <Script
        src="/charting_library-master/datafeeds/udf/dist/bundle.js"
        strategy="afterInteractive"
      />
      <Script
        src="/charting_library-master/charting_library/charting_library.standalone.js"
        strategy="afterInteractive"
        onLoad={() => {
          setScriptsLoaded(true);
        }}
      />
      {scriptsLoaded && (
        <div
          ref={widgetRef}
          style={{ border: "1px solid red ", width: "100%", height: "600px" }}
        />
      )}
    </>
  );
};

export default TradingViewWidget;
