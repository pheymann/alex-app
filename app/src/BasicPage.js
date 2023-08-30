import Header from "./Header";

export default function BasicPage(props) {
  return (
    <div>
      <Header awsContext={props.awsContext}/>

      {props.children}
    </div>
  );
}
