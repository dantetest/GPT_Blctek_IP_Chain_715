const foundationItems = [
  "不可变 Dataset Version 与 Manifest",
  "支付宝确认交付后结算与分账",
  "私有加密 P2P 交付",
  "退款、争议、存证与审计闭环",
];

export default function HomePage() {
  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">BlctekIP IP-Chain</p>
        <h1>可信数据版本，受控交易交付</h1>
        <p className="lead">
          面向 AI 训练数据的登记、存证、合规辅助审查和确认交付后结算平台。
        </p>
        <div className="status">工程初始化中 · Iteration 1</div>
      </section>
      <section className="grid" aria-label="产品核心能力">
        {foundationItems.map((item) => (
          <article key={item} className="card">
            <h2>{item}</h2>
            <p>该能力将在独立迭代中通过后端、前端、测试和审计一起交付。</p>
          </article>
        ))}
      </section>
    </main>
  );
}
