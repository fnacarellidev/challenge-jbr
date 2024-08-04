import { useNavigate } from "react-router-dom";
import React, { useState, useEffect } from "react"

export default function SearchBar() {

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
		if (cnj === "" || courtOfOrigin === "") {
			setErrMsg("⚠️  Ambos os campos devem estar preenchidos para realizar a consulta.")
			setErr(true)
			return
		}

		const endpoint = "http://localhost:8080/graphql"
		const query = {
			query: "query($cnj: String!, $court_of_origin: String!) { court_case(cnj: $cnj, court_of_origin: $court_of_origin) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
			variables: {
				cnj: cnj,
				court_of_origin: courtOfOrigin
			}
		}
		const res = await fetch(endpoint, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(query),
		})

		const graphql = await res.json()
		if ('errors' in graphql) {
			setErrMsg(`⚠️  Não foi encontrado nenhum processo de cnj ${cnj} no tribunal ${courtOfOrigin}`)
			setErr(true)
		}
		else {
			setCourtCase(graphql.data)
			setErr(false)
		}
	}

	return (
		<>
			{ err &&
				<p className="search-page-error-msg">{errMsg}</p>
			}
			<div className="search-bar-wrapper">
				<input className="search-page-input" style={{ width: '150px' }} type="text" onChange={(e) => setCourtOfOrigin(e.target.value)} placeholder="Tribunal" />
				<input className="search-page-input cnj-input" type="text" onChange={(e) => setCnj(e.target.value)} placeholder="Número de processo" />
				<button className="search-page-button" onClick={fetchCourtCase}>Buscar</button>
			</div>
		</>
	)
}
