import Link from "next/link";

const clues = [
  { number: "01", title: "Изучайте дело", text: "Читайте материалы, сверяйте показания и собирайте картину преступления по частям." },
  { number: "02", title: "Ведите допрос", text: "Задавайте собственные вопросы. Каждый подозреваемый отвечает в рамках своей истории и характера." },
  { number: "03", title: "Назовите убийцу", text: "Проверьте версии экспертизами и сдайте итоговый отчёт, когда будете уверены в выводах." },
];

export default function Home() {
  return (
    <main>
      <header className="site-header">
        <Link className="wordmark" href="/" aria-label="ДетектИИв — главная">ДЕТЕКТ<span>ИИ</span>В</Link>
        <nav aria-label="Основная навигация"><a href="#game">Об игре</a><Link href="/support">Поддержка</Link></nav>
      </header>
      <section className="hero" id="game">
        <div className="case-label"><span>Дело № 001</span><span>Доступ разрешён</span></div>
        <div className="hero-copy">
          <p className="eyebrow">Интерактивная детективная игра</p>
          <h1>Истина не лежит<br />на поверхности.</h1>
          <p className="lead">Допрашивайте подозреваемых, сопоставляйте улики и раскройте преступление, которое создаётся специально для вас.</p>
          <div className="hero-actions">
            <span className="store-button disabled" aria-label="Скоро в App Store"> <b>App Store</b><small>Скоро</small></span>
            <span className="store-button disabled" aria-label="Скоро в Google Play">▶ <b>Google Play</b><small>Скоро</small></span>
          </div>
        </div>
        <div className="evidence-board" aria-hidden="true">
          <div className="photo suspect"><span>ПОДОЗРЕВАЕМЫЙ</span></div><div className="photo evidence"><span>УЛИКА 04</span></div>
          <div className="thread thread-a" /><div className="thread thread-b" /><div className="stamp">СЕКРЕТНО</div>
          <p className="note">Кто лжёт?<br />Проверить алиби.</p>
        </div>
        <div className="scroll-mark">Листайте дело <span>↓</span></div>
      </section>
      <section className="investigation">
        <div className="section-heading"><p className="eyebrow">Ход расследования</p><h2>Ваше дело.<br />Ваши вопросы.<br />Ваш вердикт.</h2></div>
        <div className="clue-list">{clues.map((clue) => <article key={clue.number}><span>{clue.number}</span><div><h3>{clue.title}</h3><p>{clue.text}</p></div></article>)}</div>
      </section>
      <footer>
        <div className="wordmark">ДЕТЕКТ<span>ИИ</span>В</div>
        <div className="footer-links"><Link href="/privacy">Конфиденциальность</Link><Link href="/terms">Условия использования</Link><Link href="/support">Поддержка</Link><Link href="/data-deletion">Удаление данных</Link></div>
        <p>© 2026 artcodefun</p>
      </footer>
    </main>
  );
}
