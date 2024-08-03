import "./styles.css"

export default function SearchPage() {

	async function fetchStuff() {
		const endpoint = "http://localhost:8080/graphql"
		const query = `{
			"query": "query($cnj: String!) { court_case(cnj: $cnj) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
				"variables":{
					"cnj": "5001682-88.2024.8.13.0672"
				}
			}`
		const res = await fetch(endpoint, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: query,
		})

		console.log(res.status)
		const json = await res.json()
		console.log(json)
	}

	return (
		<div className="search-page-wrapper">
			<h1 className="search-page-title">Buscar</h1>
			<h3>Busque um processo a partir do número unificado</h3>
			<div className="search-bar-wrapper">
				<input className="search-page-input" type="text" placeholder="Número de processo" />
				<button className="search-page-button" onClick={fetchStuff}>Buscar</button>
			</div>
		</div>
	)
}
