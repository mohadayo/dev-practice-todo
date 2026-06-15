export const metadata = {
  title: "Dev Practice TODO",
  description: "GitHub Flow 練習用の最小 TODO アプリ",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body
        style={{
          fontFamily: "system-ui, sans-serif",
          maxWidth: 640,
          margin: "40px auto",
          padding: "0 16px",
        }}
      >
        {children}
      </body>
    </html>
  );
}
