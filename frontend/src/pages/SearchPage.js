import "./styles.css"

export default function SearchPage() {
	return (
		<div className="search-page-wrapper">
			<h1 className="search-page-title">Buscar</h1>
			<h3>Busque um processo a partir do número unificado</h3>
			<div className="search-bar-wrapper">
				<input className="search-page-input" type="text" placeholder="Número de processo" />
				<button className="search-page-button" onClick={() => console.log('hello')}>Buscar</button>
			</div>
		</div>
	)
}
