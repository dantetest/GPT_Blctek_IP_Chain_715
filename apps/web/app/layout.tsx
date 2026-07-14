import type { Metadata } from "next";
import "./styles.css";

export const metadata: Metadata = {
  title: "BlctekIP IP-Chain",
  description: "AI训练数据登记存证、合规辅助审查与受控交易平台",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
