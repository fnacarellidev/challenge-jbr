import "./styles.css"

export default function SearchPage() {
	return (
		<div class="search-page-wrapper">
			<h1 class="search-page-title">Buscar</h1>
			<h3>Busque um processo a partir do número unificado</h3>
			<div class="search-bar-wrapper">
				<input class="search-page-input" type="text" placeholder="Número de processo" />
				<button class="search-page-button">Buscar</button>
			</div>
		</div>
	)
}
