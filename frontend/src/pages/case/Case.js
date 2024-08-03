import { useLocation } from "react-router-dom";

export default function Case() {
	const { state } = useLocation()
	const { courtCase } = state || {}

	console.log(courtCase)

	return (
		<>
			<h1> { courtCase.court_case.plaintiff } VS { courtCase.court_case.defendant } </h1>
		</>
	)
}
