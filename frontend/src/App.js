import { BrowserRouter, Routes, Route } from "react-router-dom";
import SearchPage from "./pages/SearchPage"

export default function App() {
	return (
		<>
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<SearchPage />} />
					<Route path="/case" element={<h1> Hello World </h1>} />
				</Routes>
			</BrowserRouter>
		</>
	);
}
