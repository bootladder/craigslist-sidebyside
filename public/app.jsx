class App extends React.Component {
    constructor(props) {
      super(props);
      this.addButtonClicked = this.addButtonClicked.bind(this);
      this.doDeleteRequest = this.doDeleteRequest.bind(this);
      this.getCraigslistUrlSet = this.getCraigslistUrlSet.bind(this);
    }
    state = {
        urls: [],
        children: <td>hello</td>
    }
    componentDidMount(){
        console.log("app mounted")
        this.getCraigslistUrlSet(0)
    }

    getCraigslistUrlSet(setIndex){
        console.log("getCraigslistUrlSet: index " + setIndex)
        fetch("http://localhost:8080/api/" + setIndex , {
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
            globals.key = globals.key + 1
            children.push(
<td key={j}>
    <div>
        <CraigslistQueryColumn 
            url={urls[j]}
            columnIndex={j}
            doDeleteRequest={this.doDeleteRequest}
            globalKey={globals.key}
            />
    </div>
</td> 
            )
        }
        return children
    }

    addButtonClicked(setIndex) {
        console.log("Add Button Clicked")

        var myJsonRequestObj = {
            setIndex: setIndex
        };

        fetch("http://localhost:8080/api/" , {
            method: "PUT",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(myJsonRequestObj)
        })
        .then(response => response.json())
        .then(data =>
            this.updateUrls(data.urls)
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    doDeleteRequest(index) {
        console.log("Delete Button Clicked")

        var myJsonRequestObj = {
            columnIndex: index
        };
        fetch("http://localhost:8080/api/" , {
            method: "DELETE",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(myJsonRequestObj)
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
                <Button onClick={()=>this.addButtonClicked(globals.currentSetIndex)} size="small" color="primary">Add New Column</Button>
                <Button onClick={()=>this.getCraigslistUrlSet(1)} size="small" color="primary">Load URL Set #1</Button>
                <Button onClick={()=>this.getCraigslistUrlSet(2)} size="small" color="primary">Load URL Set #2</Button>
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
