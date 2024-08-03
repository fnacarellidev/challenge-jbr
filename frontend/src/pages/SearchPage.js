import { useNavigate } from "react-router-dom";
import React, { useState, useEffect } from "react"
import "./styles.css"

export default function SearchPage() {

	const [cnj, setCnj] = useState("")
	const [err, setErr] = useState(false)
	const [errMsg, setErrMsg] = useState("")
	const navigate = useNavigate();

	useEffect(() => {
		if (!err && cnj) navigate("/case")
	}, [err])

	async function fetchCourtCase() {
		const endpoint = "http://localhost:8080/graphql"
		const query = `{
			"query": "query($cnj: String!) { court_case(cnj: $cnj) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
				"variables":{
					"cnj": "${cnj}"
				}
			}`
		const res = await fetch(endpoint, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: query,
		})

		const graphql = await res.json()
		setErr('errors' in graphql)
		if ('errors' in graphql)
			setErrMsg(`⚠️ Não foi encontrado nenhum processo de cnj ${cnj}`)
	}

	return (
		<div className="search-page-wrapper">
			<h1 className="search-page-title">Buscar</h1>
			<h3>Busque um processo a partir do número unificado</h3>
			{ err &&
				<p className="search-page-error-msg">{errMsg}</p>
			}
			<div className="search-bar-wrapper">
				<input className="search-page-input" type="text" onChange={(e) => setCnj(e.target.value)} placeholder="Número de processo" />
				<button className="search-page-button" onClick={fetchCourtCase}>Buscar</button>
			</div>
		</div>
	)
}
