import { useNavigate } from "react-router-dom";
import React, { useState, useEffect } from "react"
import "./styles.css"

export default function SearchPage() {

	const [cnj, setCnj] = useState("")
	const [courtOfOrigin, setCourtOfOrigin] = useState("")
	const [err, setErr] = useState(false)
	const [errMsg, setErrMsg] = useState("")
	const [courtCase, setCourtCase] = useState(null)
	const navigate = useNavigate();

	useEffect(() => {
		if (!err && courtCase) navigate("/case", { state: { courtCase } })
	}, [err, courtCase])

	async function fetchCourtCase() {
		const endpoint = "http://localhost:8080/graphql"
		const query = `{
			"query": "query($cnj: String!, $court_of_origin: String!) { court_case(cnj: $cnj, court_of_origin: $court_of_origin) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
				"variables":{
					"cnj": "${cnj}",
					"court_of_origin": "${courtOfOrigin}"
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
		if ('errors' in graphql) {
			setErrMsg(`⚠️ Não foi encontrado nenhum processo de cnj ${cnj}`)
			setErr(true)
		}
		else {
			setCourtCase(graphql.data)
			setErr(false)
		}
	}

	return (
		<div className="search-page-wrapper">
			<h1 className="search-page-title">Buscar</h1>
			<h3 style={{ textAlign: 'center', fontWeight: '400'}}>Busque um processo a partir do número unificado</h3>
			{ err &&
				<p className="search-page-error-msg">{errMsg}</p>
			}
			<div className="search-bar-wrapper">
				<input className="search-page-input" style={{ width: '150px' }} type="text" onChange={(e) => setCourtOfOrigin(e.target.value)} placeholder="Tribunal" />
				<input className="search-page-input cnj-input" type="text" onChange={(e) => setCnj(e.target.value)} placeholder="Número de processo" />
				<button className="search-page-button" onClick={fetchCourtCase}>Buscar</button>
			</div>
		</div>
	)
}
