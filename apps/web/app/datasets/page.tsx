import Link from "next/link";
import { listDatasets } from "../../lib/datasets";

export default async function DatasetsPage() {
  const result = await listDatasets();

  return (
    <main className="shell dashboard-shell">
      <header className="page-header">
        <div>
          <p className="eyebrow">Seller workspace</p>
          <h1 className="page-title">数据集管理</h1>
          <p className="lead compact-lead">创建数据集、管理不可变版本，并跟踪 Manifest 与合规材料状态。</p>
        </div>
        <Link className="primary-action" href="/datasets/new">创建数据集</Link>
      </header>

      {!result.ok ? <div className="notice warning-notice">{result.message}</div> : null}

      <section className="dataset-grid" aria-label="数据集列表">
        {result.datasets.map((dataset) => (
          <article className="dataset-card" key={dataset.id}>
            <div className="dataset-card-topline">
              <span className="status-badge">{dataset.status}</span>
              <span className="revision">Revision {dataset.revision}</span>
            </div>
            <h2>{dataset.title}</h2>
            <p className="dataset-slug">/{dataset.slug}</p>
            <p>{dataset.description || "尚未填写数据集说明。"}</p>
            <div className="dataset-meta">
              <span>所有者：{dataset.ownerType}</span>
              <span>更新：{new Date(dataset.updatedAt).toLocaleString("zh-CN")}</span>
            </div>
          </article>
        ))}
      </section>

      {result.ok && result.datasets.length === 0 ? (
        <section className="empty-state">
          <p className="eyebrow">No datasets</p>
          <h2>还没有数据集</h2>
          <p>先创建数据集，再通过 Data Agent 生成并绑定 Manifest。</p>
          <Link className="primary-action" href="/datasets/new">创建第一个数据集</Link>
        </section>
      ) : null}
    </main>
  );
}
