import { createFileRoute } from '@tanstack/react-router'
import ImageGenerator from '../components/ImageGenerator'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="p-2">
      <h1 className="text-3xl font-bold mb-4">Go Gopher LGTM Generator</h1>
      <ImageGenerator />
    </div>
  )
}
