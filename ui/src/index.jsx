import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { StandingsTable } from './StandingsTable';
import { GamesTable } from './GamesTable';
import { Navigation } from './Nav';
import PlayerForm  from './form/PlayerForm';
import GameForm from './form/GameForm';
import {
  HashRouter,
  Routes,
  Route
} from "react-router-dom";
import * as serviceWorker from './serviceWorker';
import 'bootstrap/dist/css/bootstrap.min.css';
import { Container } from 'react-bootstrap';

ReactDOM.render(
  <HashRouter>
    <Navigation />
    <Container className="undernav">
      <Routes>
        <Route path="/standings" element={<StandingsTable/>} />
        <Route path="/games" element={<GamesTable/>} />
        <Route path="/new/game" element={<GameForm/>} />
        <Route path="/new/player" element={<PlayerForm/>} />
        <Route path="/" element={<App/>} />
      </Routes>
    </Container>
  </HashRouter>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
