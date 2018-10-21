
class App extends React.Component {
    constructor(props) {
      super(props);
      this.doRequest = this.doRequest.bind(this);
      this.getCraigslistUrls = this.getCraigslistUrls.bind(this);
    }
    state = {
        dummystate: "initial dummy state"
    }
    componentDidMount(){
        this.getCraigslistUrls()
    }
    getCraigslistUrls(){
        console.log("get craigslistURLS")

        fetch("http://localhost:8080/api/" , {
            method: "GET"
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                craigslistUrlsResponse: data.urls,
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }
    saveColumnInfo(e) {
        console.log("App saveColumnInfo")
    }
    doRequest(craigslistSearchURL){
        console.log("app do request" + craigslistSearchURL)

        var myJsonRequestObj = {
            searchURL: encodeURIComponent(craigslistSearchURL),
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myJsonRequestObj)
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
        console.log("rendering app with state " + this.state)
        console.log("rendering app with urls " + JSON.stringify(this.state.craigslistUrlsResponse))
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
        url="https://newyork.craigslist.org/d/architect-engineer-cad/search/egr"
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

