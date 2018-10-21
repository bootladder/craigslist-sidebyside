
class App extends React.Component {
    constructor(props) {
      super(props);
      this.doRequest = this.doRequest.bind(this);
    }
    state = {
        dummystate: "initial dummy state"
    }
    saveColumnInfo(e) {
        console.log("App saveColumnInfo")
    }
    doRequest(craigslistSearchURL){
        console.log("app do request" + craigslistSearchURL)
        this.setState({
                dummystate: "changeddummystate"
        })
        const myURL="https://newyork.craigslist.org/d/architect-engineer-cad/search/egr"

        var myObj = {
            searchURL: encodeURIComponent(craigslistSearchURL),
            hardsearchURL: encodeURIComponent(myURL),
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myObj)
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                users: data,
                craigslistQueryResponse: data.response,
                message: "ok!!",
                dummystate: "ok finally changed state up here!!",
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }
    render() {
        console.log("rendering app with state " + this.state.dummystate)
        return (
            <div>
                <h1>This is the App</h1>
                <hr/>
                <div className="search-table-outter wrapper">
                <table className="search-table inner">
            <tbody>
        <tr>
    <td>
<div>
    <CraigslistQueryColumn 
        saveColumnInfo={this.saveColumnInfo} 
        doRequest={this.doRequest} 
        queryResponseData={this.state.craigslistQueryResponse} 
        />
</div>
    </td> 
        </tr>
            </tbody>
                </table>
                </div>
            </div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('app'));

