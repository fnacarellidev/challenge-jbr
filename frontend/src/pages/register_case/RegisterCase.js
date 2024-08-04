import { useState } from "react"
import { useNavigate } from "react-router-dom";
import "./styles.css"

export default function RegisterCase() {
	const [cnj, setCnj] = useState("")
	const [err, setErr] = useState(false)
	const [errMsg, setErrMsg] = useState("")
	const [plaintiff, setPlaintiff] = useState("")
	const [defendant, setDefendant] = useState("")
	const [startDate, setStartDate] = useState("")
	const [courtOfOrigin, setCourtOfOrigin] = useState("")
	const [updates, setUpdates] = useState([]);
	const navigate = useNavigate()

	function handleChange(e, index) {
		const { name, value } = e.target;
		const list = [...updates];
		list[index][name] = value;
		setUpdates(list);
	};

	function handleAddUpdate() {
		setUpdates([...updates, { update_date: '', update_details: '' }]);
	};

	function handleRemoveUpdate(index) {
		const list = [...updates];
		list.splice(index, 1);
		setUpdates(list);
	};

	async function handleSubmit(e) {
		e.preventDefault();
		const endpoint = "http://localhost:8080/graphql"
		const copyUpdates = [...updates]
		copyUpdates.forEach((update) => {
			update.update_date = new Date(update.update_date).toISOString() // Need to do this copy, when I try to do that directly on the reactive variable it messes with the date input
		})
		const query = {
			query: "mutation new_court_case($cnj: String!, $plaintiff: String!, $defendant: String!, $court_of_origin: String!, $start_date: DateTime!, $updates: [CaseUpdateInput]) { new_court_case(cnj: $cnj, plaintiff: $plaintiff, defendant: $defendant, court_of_origin: $court_of_origin, start_date: $start_date, updates: $updates) { cnj plaintiff defendant court_of_origin start_date updates { update_date update_details } } }",
			variables: {
				cnj: cnj,
				plaintiff: plaintiff,
				defendant: defendant,
				court_of_origin: courtOfOrigin,
				start_date: new Date(startDate).toISOString(),
				updates: copyUpdates
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
			setErr(true)
			setErrMsg("error: " + graphql.errors[0].message)
		}
		else {
			const courtCase = {
				court_case: graphql.data.new_court_case
			}
			navigate("/case", { state: { courtCase } })
		}
	};

	return (
		<form onSubmit={handleSubmit}>
			<input type="text" className="form-input pretty-input" name="cnj" placeholder="CNJ" value={cnj} onChange={(e) => setCnj(e.target.value)} />
			<input type="text" className="form-input pretty-input" name="plaintiff" placeholder="Autor" value={plaintiff} onChange={(e) => setPlaintiff(e.target.value)} />
			<input type="text" className="form-input pretty-input" name="defendant" placeholder="Réu" value={defendant} onChange={(e) => setDefendant(e.target.value)} />
			<input type="text" className="form-input pretty-input" name="court_of_origin" placeholder="Tribunal" value={courtOfOrigin} onChange={(e) => setCourtOfOrigin(e.target.value)} />
			<input type="date" className="form-input pretty-input" name="start_date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />

			{ updates.map((update, index) => (
				<div key={index}>
				<input type="datetime-local" className="pretty-input" name="update_date" value={update.update_date} onChange={e => handleChange(e, index)} />
				<input type="text" className="pretty-input" name="update_details" placeholder="Update Details" value={update.update_details} onChange={e => handleChange(e, index)} />
					{ updates.length > 0 && (
						<button type="button" className="pretty-btn" onClick={() => handleRemoveUpdate(index)}>Remove</button>
					)}
				</div>
			))}

			{ err &&
				<p className="search-page-error-msg">{errMsg}</p>
			}

			<button type="button" className="pretty-btn" style={{ display: 'block', marginBottom: '4px' }} onClick={handleAddUpdate}>Adicionar movimentação</button>
			<button type="submit" className="pretty-btn">Submit</button>
		</form>
	);

}
