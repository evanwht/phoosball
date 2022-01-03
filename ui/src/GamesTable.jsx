import React, { Component } from 'react';
import Table from 'react-bootstrap/Table';
import API from './api/api';
import GameEditButton from './form/GameEditButton.jsx';

export class GamesTable extends Component {
    constructor(props) {
        super(props)
        this.state = {
            games: [
                {
                    id: 1,
                    played: {
                        epochSecond: "07-04-2021"
                    },
                    team1: {
                        defense: {
                            id: 1,
                            name: 'Evan White'
                        },
                        offense: {
                            id: 2,
                            name: 'Thomas Mckenna'
                        }
                    },
                    team2: {
                        defense: {
                            id: 3,
                            name: 'Zach Volz'
                        },
                        offense: {
                            id: 4,
                            name: 'Manny Shahugan'
                        }
                    },
                    team1Half: 5,
                    team2Half: 3,
                    team1Final: 10,
                    team2Final: 6
                }
            ]
        }
        this.refreshPage = this.refreshPage.bind(this);
    }

    componentDidMount() {
        API.get('games')
            .then(res => {
                const games = res.data;
                this.setState({ games });
            });
    }

    refreshPage(refresh) {
        if (refresh) {
            this.componentDidMount();
        }
    }


    compareGames(g1, g2) {
        if (g1.played.epochSecond < g2.played.epochSecond) {
            return 1;
        } else if (g1.played.epochSecond > g2.played.epochSecond) {
            return -1;
        }
        return 0;
    }

    renderTableData() {
        let games = this.state.games;
        games.sort(this.compareGames);
        return this.state.games.map((game, i) => {
            const { id, played, team1, team2, team1Half, team2Half, team1Final, team2Final } = game;
            const team_1 = team1.defense.name + " - " + team1.offense.name;
            const team_2 = team2.defense.name + " - " + team2.offense.name;
            const half_scores = team1Half + " - " + team2Half;
            const final_scores = team1Final + " - " + team2Final;
            const winners = team1Final > team2Final ? team_1 : team_2;
            const losers = team1Final > team2Final ? team_2 : team_1;
            return (
                <tr key={id}>
                    <td className="align-middle text-center">{new Date(played).toLocaleString()}</td>
                    <td className="align-middle winners text-center">{winners}</td>
                    <td className="align-middle losers text-center">{losers}</td>
                    <td className="align-middle text-center">{half_scores}</td>
                    <td className="align-middle text-center">{final_scores}</td>
                    <td className="align-middle text-center">
                        <GameEditButton
                            id={id}
                            refresh={this.refreshPage}
                            played={played}
                            p1={team1.defense.id}
                            p2={team1.offense.id}
                            p3={team2.defense.id}
                            p4={team2.offense.id}
                            t1h={team1Half}
                            t2h={team2Half}
                            t1f={team1Final}
                            t2f={team2Final}
                        />
                    </td>
                </tr>
            )
        })
    }

    render() {
        return (
            <Table className="rounded table-borderless" size="sm" striped hover variant="dark">
                <thead>
                    <tr>
                        <th className="text-center">Date</th>
                        <th className="text-center">Winners</th>
                        <th className="text-center">Losers</th>
                        <th className="text-center">Half</th>
                        <th className="text-center">Final</th>
                        <th className="text-center">Edit</th>
                    </tr>
                </thead>
                <tbody>
                    {this.renderTableData()}
                </tbody>
            </Table>
        );
    }
}
