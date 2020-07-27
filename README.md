# emergent reboot in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/emer/emergent)](https://goreportcard.com/report/github.com/emer/emergent)
[![GoDoc](https://godoc.org/github.com/emer/emergent?status.svg)](https://godoc.org/github.com/emer/emergent)
[![Travis](https://travis-ci.com/emer/emergent.svg?branch=master)](https://travis-ci.com/emer/emergent)

This is the new home of the *emergent* neural network simulation software, developed primarily by the CCN lab at the University of Colorado Boulder (now at UC Davis): https://ccnlab.org.  We have decided to completely reboot the entire enterprise from the ground up, with a much more open, general-purpose design and approach.

See [Wiki Install](https://github.com/emer/emergent/wiki/Install) for installation instructions (note: Go 1.13 and newer are now required!), and the [Wiki Rationale](https://github.com/emer/emergent/wiki/Rationale) and [History](https://github.com/emer/emergent/wiki/History) pages for a more detailed rationale for the new version of emergent, and a history of emergent (and its predecessors).

See the [ra25 example](https://github.com/emer/leabra/blob/master/examples/ra25/README.md) in the `leabra` package for a complete working example (intended to be a good starting point for creating your own models), and any of the 26 models in the [Comp Cog Neuro sims](https://github.com/CompCogNeuro/sims) repository which also provide good starting points.  See the [etable wiki](https://github.com/emer/etable/wiki) for docs and example code for the widely-used etable data table structure, and the `family_trees` example in the CCN textbook sims which has good examples of many standard network representation analysis techniques (PCA, cluster plots, RSA).

# Current Status / News

* April 2020: Version 1.0 of GoGi GUI is now released, and we have updated all module dependencies accordingly. *We now recommend using the go modules instead of GOPATH* -- the [Wiki Install](https://github.com/emer/emergent/wiki/Install) instructions have been updated accordingly.

* 12/30/2019: Version 1.0.0 released!  The [Comp Cog Neuro sims](https://github.com/CompCogNeuro/sims) that accompany the [CCN Textbook](https://github.com/CompCogNeuro/ed4) are now complete and have driven extensive testing and bugfixing.

* 3/2019: Python interface is up and running!  See the `python` directory in `leabra` for the [README](https://github.com/emer/leabra/blob/master/python/README.md) status and how to give it a try.  You can run the full `leabra/examples/ra25` code using Python, including the GUI etc.

* 2/2019: Initial implementation and benchmarking (see `examples/bench` for details -- shows that the Go version is comparable in speed to C++).

# Key Features

* Currently focused exclusively on implementing the biologically-based `Leabra` algorithm (now in a separate repository), which is not at all suited to implementation in current popular neural network frameworks such as `PyTorch`.  Leabra uses point-neurons and competitive inhibition, and has sparse activity levels and ubiquitous fully recurrent bidirectional processing, which enable / require novel optimizations for how simulated neurons communicate, etc.

* Go-based code can be compiled to run entire models.  Instead of creating and running everything in the *emergent* GUI, the process is much more similar to how e.g., PyTorch and other current frameworks work.  You write code to configure your model, and directly call functions that run your model, etc.  This gives you full, direct, transparent control over everything that happens in your model, as opposed to the previous highly opaque nature of [C++ emergent](https://github.com/emer/cemer).

* Although we will be updating our core library (`package` in Go) code with bug fixes, performance improvements, and new algorithms, we encourage users who have invested in developing a particular model to fork their own copy of the codebase and use that to maintain control over everything.  Once we make our official release of the code, the raw algorithm code is essentially guaranteed to remain fairly stable and encapsulated, so further changes should be relatively minimal, but nevertheless, it would be good to have an insurance policy!  The code is very compact and having your own fork should be very manageable.

* The `emergent` repository will host additional Go packages that provide support for models.  These are all designed to be usable as independently and optionally as possible.  Users running Leabra from Python for example will likely rely on relevant tools in that ecosystem instead.  An overview of some of those packages is provided below.

* We are committed to making the system fully usable from within Python, given the extensive base of Python users.  See the [leabra python README](https://github.com/emer/leabra/blob/master/python/README.md).  This includes interoperating with [PsyNeuLink](https://princetonuniversity.github.io/PsyNeuLink/) to make Leabra models accessible in that framework, and vice-versa.  Furthermore, interactive, IDE-level tools such as `Jupyter` and `nteract` can be used to interactively develop and analyze the models, etc. 

* We are leveraging the [GoGi Gui](https://github.com/goki/gi) to provide interactive 2D and 3D GUI interfaces to models, capturing the essential functionality of the original C++ emergent interface, but in a much more a-la-carte fashion.  We will also support the [GoNum](https://github.com/gonum) framework for analyzing and plotting results within Go.

# Design

* In general, *emergent* works by compiling programs into executables which you then run like any other executable. This is very different from the C++ version of emergent which was a single monolithic program attempting to have all functionality built-in. Instead, the new model is the more prevalent approach of writing more specific code to achieve more specific goals, which is more flexible and allows individuals to be more in control of their own destiny..
    + To make your own simulations, start with e.g., the `examples/ra25/ra25.go` code (or that of a more appropriate example) and copy that to your own repository, and edit accordingly.

* The `emergent` repository contains a collection of packages supporting the implementation of biologically-based neural networks.  The main package is `emer` which specifies a minimal abstract interface for a neural network.  The `etable` `etable.Table` data structure (DataTable in C++) is in a separate repository under the overall `emer` project umbrella, as are specific algorithms such as `leabra` which implement the `emer` interface.

* Go uses `interfaces` to represent abstract collections of functionality (i.e., sets of methods).  The `emer` package provides a set of interfaces for each structural level (e.g., `emer.Layer` etc) -- any given specific layer must implement all of these methods, and the structural containers (e.g., the list of layers in a network) are lists of these interfaces.  An interface is implicitly a *pointer* to an actual concrete object that implements the interface.  Thus, we typically need to convert this interface into the pointer to the actual concrete type, as in:

```Go
func (nt *Network) InitActs() {
	for _, ly := range nt.Layers {
		if ly.IsOff() {
			continue
		}
		ly.(*Layer).InitActs() // ly is the emer.Layer interface -- (*Layer) converts to leabra.Layer
	}
}
```

* The emer interfaces are designed to support generic access to network state, e.g., for the 3D network viewer, but specifically avoid anything algorithmic.  Thus, they should allow viewing of any kind of network, including PyTorch backprop nets.

* There are 3 main levels of structure: `Network`, `Layer` and `Prjn` (projection).  The network calls methods on its Layers, and Layers iterate over both `Neuron` data structures (which have only a minimal set of methods) and the `Prjn`s, to implement the relevant computations.  The `Prjn` fully manages everything about a projection of connectivity between two layers, including the full list of `Syanpse` elements in the connection.  There is no "ConGroup" or "ConState" level as was used in C++, which greatly simplifies many things.  The Layer also has a set of `Pool` elements, one for each level at which inhibition is computed (there is always one for the Layer, and then optionally one for each Sub-Pool of units (*Pool* is the new simpler term for "Unit Group" from C++ emergent).

* Layers have a `Shape` property, using the `etensor.Shape` type (see `etable` package), which specifies their n-dimensional (tensor) shape.  Standard layers are expected to use a 2D Y*X shape (note: dimension order is now outer-to-inner or *RowMajor* now), and a 4D shape then enables `Pools` ("unit groups") as hypercolumn-like structures within a layer that can have their own local level of inihbition, and are also used extensively for organizing patterns of connectivity.

# Packages

Here are some of the additional supporting packages:

* `emer` *only* has the primary abstract Network interfaces (previously had put other random things in there, but the new policy is to keep everything in separate packages, as that seems to be where things end up eventually as they are better developed).

* `params` has the parameter-styling infrastructure (e.g., `params.Set`, `params.Sheet`, `params.Sel`), which implement a powerful, flexible, and efficient CSS style-sheet approach to parameters.  See the [Wiki Params](https://github.com/emer/emergent/wiki/Params) page for more info.

* `env` has an interface for environments, which encapsulates all the counters and timing information for patterns that are presented to the network, and enables more of a mix-and-match ability for using different environments with different networks.  See [Wiki Env](https://github.com/emer/emergent/wiki/Env) page for more info.

* `netview` provides the `NetView` interactive 3D network viewer, implemented in the GoGi 3D framework.

* `prjn` is a separate package for defining patterns of connectivity between layers (i.e., the `ProjectionSpec`s from C++ emergent).  This is done using a fully independent structure that *only* knows about the shapes of the two layers, and it returns a fully general bitmap representation of the pattern of connectivity between them.  The `leabra.Prjn` code then uses these patterns to do all the nitty-gritty of connecting up neurons.  This makes the projection code *much* simpler compared to the ProjectionSpec in C++ emergent, which was involved in both creating the pattern and also all the complexity of setting up the actual connections themselves.  This should be the *last* time any of those projection patterns need to be written (having re-written this code too many times in the C++ version as the details of memory allocations changed).

* `patgen` supports various general-purpose pattern-generation algorithms, as implemented in `taDataGen` in C++ emergent (e.g., `PermutedBinary` and `FlipBits`).

* `esg` is the *emergent stochastic / sentence generator* -- parses simple grammars that generate random events (sentences) -- can be a good starting point for generating more complex environments.

* `popcode` supports the encoding and decoding of population codes -- distributed representations of numeric quantities across a population of neurons.  This is the `ScalarVal` functionality from C++ emergent, but now completely independent of any specific algorithm so it can be used anywhere.

* `erand` has misc random-number generation support functionality, including `erand.RndParams` for parameterizing the type of random noise to add to a model, and easier support for making permuted random lists, etc.

* `timer` is a simple interval timing struct, used for benchmarking / profiling etc.

* `python` contains a template `Makefile` that uses [GoPy](https://github.com/goki/gopy) to generate python bindings to the entire emergent system.  See the `leabra` package version to actually run an example.

* The [etable](https://github.com/emer/etable) repository holds all of the more general-purpose "DataTable" or DataFrame (`etable.Table`) related code, which is our version of something like `pandas` or `xarray` in Python.  This includes the `etensor` n-dimensional array, `eplot` for interactive plotting of data, and basic utility packages like `minmax` and `bitslice`, and lots of data analysis tools like similarity / distance matricies, PCA, cluster plots, etc.

# TODO

Last updated: April 2020

This list is not strictly in order, but roughly so..

- [ ] write converter from Go to Python

- [ ] add python example code for interchange between pandas, xarray, tensorflow tensor stuff and etable.Table -- right now the best is to just save as .csv and load from there (esp for pandas which doesn't have tensors) -- should be able to use arrow stuff so it would be good to look into that.

- [ ] pvlv

- [ ] GPU -- see https://github.com/gorgonia/gorgonia for existing CUDA impl -- alternatively, maybe try using opengl or vulkan directly within existing gogi/gpu framework -- would work on any GPU and seems like it wouldn't be very hard and gives full control -- https://www.khronos.org/opengl/wiki/Compute_Shader -- 4.3 min version though -- maybe better to just go to vulkan?  https://community.khronos.org/t/opencl-vs-vulkan-compute/7132/6

- [x] MPI -- see [MPI Wiki page](https://github.com/emer/emergent/wiki/DMem)

- [x] virtual environment -- [eve](https://github.com/emer/eve) is under way -- some basic actual physics in place, and basic collision detection.

- [x] finalize GoGi GUI version 1.0 release -- finally done!


