class App extends React.Component {
    constructor(props) {
      super(props);
      this.addButtonClicked = this.addButtonClicked.bind(this);
      this.doDeleteRequest = this.doDeleteRequest.bind(this);
      this.getCraigslistUrlSet = this.getCraigslistUrlSet.bind(this);
      this.doRequest = this.doRequest.bind(this);
    }
    state = {
        urls: [],
        columns: []
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
        .then(data => {
            this.updateColumns(data.urls)
            globals.currentSetIndex = setIndex
        }
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    updateColumns(urls){
        let newColumns = []
        for(let j=0;j<urls.length;j++){
            let newColumn = {url: urls[j]}
            newColumns.push(newColumn)
        }
        console.log("updateColumns (will setState): columns: "+ JSON.stringify(newColumns))
        this.setState({
            columns: newColumns
        })
    }

    makeColumns(columnInfos) {
        console.log("     makeColumns: making "+columnInfos.length+"columns")

        let columns = []
        for (let j = 0; j < columnInfos.length; j++) {
            globals.key = globals.key + 1
            columns.push(
<td key={j}>
    <div>
        <CraigslistQueryColumn 
            url={columnInfos[j].url}
            response={columnInfos[j].response}
            columnIndex={j}
            doDeleteRequest={this.doDeleteRequest}
            globalKey={globals.key}
            doRequest={this.doRequest}
            />
    </div>
</td> 
            )
        }
        return columns
    }

    addButtonClicked(setIndex) {
        console.log("Add Button Clicked with setIndex: " + setIndex)

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
            this.updateColumns(data.urls)
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    doDeleteRequest(index) {
        console.log("Delete Button Clicked, currentSetIndex: " + globals.currentSetIndex + " columnIndex: " + index)

        var myJsonRequestObj = {
            setIndex: globals.currentSetIndex,
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
            this.updateColumns(data.urls)
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
    }

    validateCraigslistURL(url){
        if(url.length < 5){
            console.log("length too short, defaulting URL length was "+url.length)
            return "https://baltimore.craigslist.org/d/architect-engineer-cad/search/egr";
        }
        else return url;
    }
    doRequest(index, craigslistSearchURL){
        console.log("do request : index: " +index + " URL: " + craigslistSearchURL )

        var myJsonRequestObj = {
            searchURL: encodeURIComponent(craigslistSearchURL),
            setIndex: globals.currentSetIndex,
            columnIndex: index
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
        .then(data =>{
            let newColumns = this.state.columns
            newColumns[index].url = craigslistSearchURL
            newColumns[index].response = data.response
            this.setState({ 
                columns: newColumns
            })

        }
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }

    render() {
        console.log("rendering app with state " + JSON.stringify(this.state, null, 2))
        console.log("                   urls " + JSON.stringify(this.state.urls))

        return (
            <div>
                <h1>This is the App</h1>
                <Button onClick={()=>this.addButtonClicked(globals.currentSetIndex)} size="small" color="primary">Add New Column</Button>
                <Button onClick={()=>this.getCraigslistUrlSet(0)} size="small" color="primary">Load URL Set #1</Button>
                <Button onClick={()=>this.getCraigslistUrlSet(1)} size="small" color="primary">Load URL Set #2</Button>
                <hr/>
                <div className="search-table-outter wrapper">
                <table className="search-table inner">
            <tbody>
        <tr>
            {this.makeColumns(this.state.columns)}
        </tr>
            </tbody>
                </table>
                </div>
            </div>
        );
    }
}
