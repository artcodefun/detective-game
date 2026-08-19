import Link from "next/link";
import type { ReactNode } from "react";
export function LegalLayout({ eyebrow, title, intro, children }: { eyebrow: string; title: string; intro: string; children: ReactNode }) {
  return <main className="legal-shell"><header className="legal-header"><Link className="wordmark" href="/">ДЕТЕКТ<span>ИИ</span>В</Link><Link href="/">← Вернуться к делу</Link></header><article className="legal-document"><div className="document-title"><p className="eyebrow">{eyebrow}</p><h1>{title}</h1><p>{intro}</p></div><div className="document-content">{children}</div></article><footer className="legal-footer"><Link href="/privacy">Конфиденциальность</Link><Link href="/terms">Условия</Link><Link href="/support">Поддержка</Link><Link href="/data-deletion">Удаление данных</Link></footer></main>;
}
