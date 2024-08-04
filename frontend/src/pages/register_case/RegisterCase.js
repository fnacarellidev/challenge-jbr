import { useState } from "react"

export default function RegisterCase() {
	const [cnj, setCnj] = useState("")
	const [plaintiff, setPlaintiff] = useState("")
	const [defendant, setDefendant] = useState("")
	const [startDate, setStartDate] = useState("")
	const [courtOfOrigin, setCourtOfOrigin] = useState("")
	const [updates, setUpdates] = useState([{ update_date: '', update_details: '' }]);

	function handleChange(e, index) {
		const { name, value } = e.target;
		console.log(name, value)
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
			query: "mutation new_court_case($cnj: String!, $plaintiff: String!, $defendant: String!, $court_of_origin: String!, $start_date: DateTime!, $updates: [CaseUpdateInput]) { new_court_case(cnj: $cnj, plaintiff: $plaintiff, defendant: $defendant, court_of_origin: $court_of_origin, start_date: $start_date, updates: $updates) { cnj } }",
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
		console.log(graphql)
	};


	return (
		<form onSubmit={handleSubmit}>
			<input type="text" name="cnj" placeholder="CNJ" value={cnj} onChange={(e) => setCnj(e.target.value)} />
			<input type="text" name="plaintiff" placeholder="Autor" value={plaintiff} onChange={(e) => setPlaintiff(e.target.value)} />
			<input type="text" name="defendant" placeholder="Réu" value={defendant} onChange={(e) => setDefendant(e.target.value)} />
			<input type="text" name="court_of_origin" placeholder="Tribunal" value={courtOfOrigin} onChange={(e) => setCourtOfOrigin(e.target.value)} />
			<input type="date" name="start_date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />

			{ updates.map((update, index) => (
				<div key={index}>
					<input type="datetime-local" name="update_date" value={update.update_date} onChange={e => handleChange(e, index)} />
					<input type="text" name="update_details" placeholder="Update Details" value={update.update_details} onChange={e => handleChange(e, index)} />
					{ updates.length > 1 && (
						<button type="button" onClick={() => handleRemoveUpdate(index)}>Remove</button>
					)}
				</div>
			))}

			<button type="button" onClick={handleAddUpdate}>Adicionar movimentação</button>
			<button type="submit">Submit</button>
		</form>
	);

}
