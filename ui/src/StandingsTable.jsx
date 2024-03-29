import React, { Component } from 'react';
import Table from 'react-bootstrap/Table';
import API from "./api/api";

export class StandingsTable extends Component {
    constructor(props) {
        super(props)
        this.state = {
            standings: [
                { id: 1, name: 'Evan White', wins: 1, losses: 0 },
                { id: 2, name: 'Zach Volz', wins: 0, losses: 1 }
            ]
        }
    }

    componentDidMount() {
        API.get('standings')
            .then(res => {
                const standings = res.data;
                this.setState({ standings });
            })
    }

    renderTableData() {
        return this.state.standings.map((standing, index) => {
            const { name, wins, losses } = standing
            let percentage = wins/ losses
            return (
                <tr key={index}>
                    <td className="align-middle">{index + 1}</td>
                    <td className="align-middle">{name}</td>
                    <td className="align-middle">{wins}</td>
                    <td className="align-middle">{losses}</td>
                    <td className="align-middle">{(Math.round(percentage * 100) / 100).toFixed(2)}</td>
                </tr>
            )
        })
    }

    render() {
        return (
            <Table className="App" size="sm" striped bordered hover variant="dark">
                <thead>
                    <tr>
                        <th>Place</th>
                        <th>Player</th>
                        <th>Wins</th>
                        <th>Losses</th>
                        <th>Percentage</th>
                    </tr>
                </thead>
                <tbody>
                    {this.renderTableData()}
                </tbody>
            </Table>
        );
    }
}
