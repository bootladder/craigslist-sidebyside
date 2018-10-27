class App extends React.Component {
    constructor(props) {
      super(props);
      this.getCraigslistUrls = this.getCraigslistUrls.bind(this);
      this.addButtonClicked = this.addButtonClicked.bind(this);
    }
    state = {
        dummystate: "initial dummy state",
        urls: [],
        children: <td>hello</td>
    }
    componentDidMount(){
        console.log("app mounted")
        this.getCraigslistUrls()
    }
    getCraigslistUrls(){
        console.log("     get craigslistURLS")

        fetch("http://localhost:8080/api/" , {
            method: "GET"
        })
        .then(response => response.json())
        .then(data =>
            this.updateUrls(data.urls)
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }
    updateUrls(urls) {
        console.log("     updateUrls: "+ urls)
        var c = this.makeChildren(urls)
        this.setState({
            urls: urls,
            children: c
        })
    }

    makeChildren(urls) {
        console.log("     makeChildren: making "+urls.length+"chilrden")

        let children = []
        for (let j = 0; j < urls.length; j++) {
            children.push(
<td key={j}>
    <div>
        <CraigslistQueryColumn 
            url={urls[j]}
            columnIndex={j}
            />
    </div>
</td> 
            )
        }
        return children
    }

    addButtonClicked(e) {
        console.log("Add Button Clicked")

        fetch("http://localhost:8080/api/" , {
            method: "PUT"
        })
        .then(response => response.json())
        .then(data =>
            this.updateUrls(data.urls)
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    render() {
        console.log("rendering app with state " + JSON.stringify(this.state))
        console.log("                   urls " + JSON.stringify(this.state.urls))

        return (
            <div>
                <h1>This is the App</h1>
                <Button onClick={this.addButtonClicked} size="small" color="primary">Add New Column</Button>
                <hr/>
                <div className="search-table-outter wrapper">
                <table className="search-table inner">
            <tbody>
        <tr>
            {this.state.children}
        </tr>
            </tbody>
                </table>
                </div>
            </div>
        );
    }
}
